const fileupload = require('express-fileupload')
const fetch = require('node-fetch')
const cors = require('cors')
const express = require('express')
const webp = require('webp-converter')
const fs = require('fs')
const app = express()
const PORT = 8081

app.use(fileupload())
app.use(cors())

app.get('/', (req, res) => {
  console.log('root requested')
  res.send(`
    <html>
      <body>
        <form ref='uploadForm' 
          id='uploadForm' 
          action='/upload' 
          method='post' 
          encType="multipart/form-data">
            <input name="url" value="asdf"/>
            <input type="file" name="file" />
            <input type='submit' value='Upload!' />
        </form>     
      </body>
    </html>
  `)
})

app.post('/upload', async (req, res) => {
  console.log(`Upload: ${req.body.url}`)
  if (!req.files || Object.keys(req.files).length === 0) {
    return res.status(400).send('No files were uploaded.')
  }

  let filename = ''
  if(req.body.url == '') {
    filename = 'pending/xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
      let r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8)
      return v.toString(16)
    })
  }
  else {
    const available = await fetch('https://api.non.io/posts/url-is-available/' + req.body.url)
    if(!await available.json()) return res.status(400).send('URL is not available')
    filename = req.body.url
  }

  console.log(req.files)
  const file = req.files.file
  console.log(file)
  let extension = file.name.match(/\.[0-9a-zA-Z]+$/)
  if(!extension) return res.status(400).send('No file extension found')
  extension = extension[0].toLowerCase()

  const tmpPath = `${__dirname}/tmp/${filename + extension}`
  file.mv(tmpPath, err => {
    if(err) {
      res.writeHead(500, {'Content-Type': 'application/json'})
      res.end(JSON.stringify({status: 'error', message: err}))
      return
    }
  })

  const path = `${__dirname}/files/${filename}.webp`
  const thumbPath = `${__dirname}/thumbnails/${filename}.webp`
  //convert the image
  webp[extension == '.gif' ? 'gwebp' : 'cwebp'](tmpPath, path, "-q 80", status => {
    if(status == 100) {
      //if the conversion is successful, generate a thumbnail
      webp[extension == '.gif' ? 'gwebp' : 'cwebp'](tmpPath, thumbPath, "-resize 192 0 -q 50", status => {
        res.writeHead(200, {'Content-Type': 'application/json'})
        res.end(JSON.stringify({status :'success', path: `${filename}.webp`}))
      })
    }
    else return res.status(500).send('Conversion error')
  })
})

app.post('/move', async (req, res) => {
  let oldUrl = `${__dirname}/files/${req.body.oldUrl}`
  let newUrl = `${__dirname}/files/${req.body.url}`
  console.log(`move from ${oldUrl} to ${newUrl} requested`)

  if(fs.existsSync(oldUrl)){
    fs.rename(oldUrl, newUrl, err => {
      if(err) return res.status(500).send('Rename error')
      return res.send({status :'success', path: `${req.body.url}`})
    })
  }
  else return res.status(400).send('Previously uploaded image not found')


})

app.listen(PORT, ()=>console.log('listening on ' + PORT))
