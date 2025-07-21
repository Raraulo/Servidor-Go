--se crea la base de datos srvr
-- se corre la aplicación para crear tablas y:
insert into prsn values(default, '1701983559', 'Guido', null, now(), null, null, 'admin', '123', 1, null, 'M', null);

--probar generación de token con:
http://127.0.0.1:4000/api/login   --> POST
Body:
{
  "login": "admin",
  "pass": "123"
}

--> 
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE3NTIyNDIwMDUsIklkIjowLCJJc3N1ZWRBdCI6MTc0OTY1MDAwNSwiUGVyZmlsSWQiOjAsIlVzZXJuYW1lIjoiYWRtaW4ifQ.ZX_eMg07sEeZSVttLrO5tDW0XSMsrCogWjbZwdJ2ZMc",
  "ok": "true",
  "uid": "0",
  "nombre": "admin",
  "perfil": 0
}

http://127.0.0.1:4000/api/user   GET  --el token es sin las comillas.
Header:
token   --> eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE3NTIyNDIwMDUsIklkIjowLCJJc3N1ZWRBdCI6MTc0OTY1MDAwNSwiUGVyZmlsSWQiOjAsIlVzZXJuYW1lIjoiYWRtaW4ifQ.ZX_eMg07sEeZSVttLrO5tDW0XSMsrCogWjbZwdJ2ZMc

--> jsopn de respuesta.

