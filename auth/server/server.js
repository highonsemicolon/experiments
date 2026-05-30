const express = require('express')
const app = express()
const { auth } = require('express-oauth2-jwt-bearer')

const port = process.env.PORT || 8080

const jwtCheck = auth({
  audience: 'pathfinder-backend',
  issuerBaseURL: 'https://pathfinder-1.us.auth0.com/',
  tokenSigningAlg: 'RS256',
})

app.get('/', function (req, res) {
  res.send('welcome')
})

app.get('/authorized', jwtCheck, function (req, res) {
  res.json({
    msg: 'you are in',
    userId: req.auth.payload.sub,
    audience: req.auth.payload.aud,
    issuer: req.auth.payload.iss,
    scope: req.auth.payload.scope,
  })
})

app.listen(port)

console.log('Running on port ', port);

