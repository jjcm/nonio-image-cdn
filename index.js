const fileupload = require('express-fileupload')
const express = require('express')
const app = express()
const PORT = 8081

app.use(
  fileupload()
)

app.get('/', (req, res) => {
  res.send(`
    <html>
      <body>
        <form ref='uploadForm' 
          id='uploadForm' 
          action='/upload' 
          method='post' 
          encType="multipart/form-data">
            <input type="file" name="myFile" />
            <input type='submit' value='Upload!' />
        </form>     
      </body>
    </html>
  `)
})

app.post('/upload', (req, res) => {
  if (!req.files || Object.keys(req.files).length === 0) {
    return res.status(400).send('No files were uploaded.')
  }
  console.log(req.files)
  const file = req.files.myFile
  const path = __dirname + '/files/' + file.name

  file.mv(path, err => {
    if(err) {
      console.error(err)
      res.writeHead(500, {'Content-Type': 'application/json'})
      res.end(JSON.stringify({status: 'error', message: err}))
      return
    }
  })

  res.writeHead(200, {'Content-Type': 'application/json'})
  res.end(JSON.stringify({status :'success', path: '/files/' + file.name}))
})

app.listen(PORT, ()=>console.log('listening on ' + PORT))