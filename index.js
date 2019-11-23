const fileupload = require('express-fileupload')
const express = require('express')
const app = express()
const PORT = 8081

app.use(
  fileupload()
)
