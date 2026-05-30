const express = require('express')
const { auth } = require('express-oauth2-jwt-bearer')
const cors = require('cors')

const port = process.env.PORT || 8081

const jwtCheck = auth({
  audience: 'pathfinder-backend',
  issuerBaseURL: 'https://pathfinder-1.us.auth0.com/',
  tokenSigningAlg: 'RS256',
})

const app = express()

app.use(cors({
  origin: 'http://localhost:8080',
  credentials: true,
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization']
}))

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

