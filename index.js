const fileupload = require('express-fileupload')
const fetch = require('node-fetch')
const cors = require('cors')
const express = require('express')
const webp = require('webp-converter')
const app = express()
const PORT = 8081

app.use(fileupload())
app.use(cors())

app.get('/', (req, res) => {
  res.send(`
    <html>
      <body>
        <form ref='uploadForm' 
          id='uploadForm' 
          action='/upload' 
          method='post' 
          encType="multipart/form-data">
            <input name="url" value="asdf"/>
            <input type="file" name="myFile" />
            <input type='submit' value='Upload!' />
        </form>     
      </body>
    </html>
  `)
})

app.post('/upload', async (req, res) => {
  if (!req.files || Object.keys(req.files).length === 0) {
    return res.status(400).send('No files were uploaded.')
  }

  let options = {}
  options.headers = {
    Authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InF3ZXIyMzE0MjM0QGFzZGZjLmNvbSIsImV4cGlyZXNBdCI6MTU3NDg0MTUzMH0.55hnoyKUj6TMzgA6vSGi7InRSL0mDHDBXtGSGHyvi4E"
  }

  const available = await fetch('https://api.non.io/posts/url-is-available/' + req.body.url, options)
  if(!available) return res.status(400).send('URL is not available')


  const file = req.files.file
  const extension = file.name.match(/\.[0-9a-z]+$/)
  if(!extension) return res.status(400).send('No file extension found')

  const tmpPath = `${__dirname}/tmp/${req.body.url + extension[0]}`
  file.mv(tmpPath, err => {
    if(err) {
      res.writeHead(500, {'Content-Type': 'application/json'})
      res.end(JSON.stringify({status: 'error', message: err}))
      return
    }
  })

  const path = `${__dirname}/files/${req.body.url}.webp`
  webp[extension == '.gif' ? 'gwebp' : 'cwebp'](tmpPath, path, "-q 80", status => {
    if(status == 100) {
      res.writeHead(200, {'Content-Type': 'application/json'})
      res.end(JSON.stringify({status :'success', path: `${req.body.url}.webp`}))
    }
    else return res.status(500).send('Conversion error')
  })

})

app.listen(PORT, ()=>console.log('listening on ' + PORT))