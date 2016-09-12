#!/usr/bin/python
import json
from bigchainwrapper import *
from http.server import BaseHTTPRequestHandler,HTTPServer
from os import curdir, sep
import cgi

PORT_NUMBER = 8081

#This class will handles any incoming request from
#the browser
class myHandler(BaseHTTPRequestHandler):

	#Handler for the GET requests
	def do_GET(self):

		if self.path=="/":
			self.path="/index_example3.html"

		try:
			#Check the file extension required and
			#set the right mime type

			sendReply = True
			if self.path.endswith(".html"):
				mimetype='text/html'
				sendReply = True
			if self.path.endswith(".jpg"):
				mimetype='image/jpg'
				sendReply = True
			if self.path.endswith(".gif"):
				mimetype='image/gif'
				sendReply = True
			if self.path.endswith(".js"):
				mimetype='application/javascript'
				sendReply = True
			if self.path.endswith(".css"):
				mimetype='text/css'
				sendReply = True

			if sendReply == True: #Open the static file requested and send it f =
				open(curdir + sep + self.path)
				self.send_response(200)
				self.send_header('Content-type',mimetype)
				self.end_headers()
				# Send message back to client
				message= "hello bigchain"
				# Write content as utf-8 data
				self.wfile.write(bytes(message, "utf8"))
				#self.wfile.write(f.read())
				f.close()
			return

		except IOError:
			self.send_error(404,'File Not Found: %s' % self.path)

	#Handler for the POST requests
	def do_POST(self):
		if self.path=="/create":

			data_string = self.rfile.read(int(self.headers['Content-Length']))

			readable = json.loads(data_string.decode('utf-8'))

			print(data_string, readable)

			user = readable['user']
			asset = readable['asset']

			self.send_response(200)
			self.end_headers()
		
			tx, result = bigchainwrapper.newAsset(user, asset)

			self.wfile.write(bytes('{"user":"' + user + '","tx":'+  str(tx) + '}',"utf8"))
			return

		if self.path=="/transfer":

			data_string = self.rfile.read(int(self.headers['Content-Length']))

			readable = json.loads(data_string.decode('utf-8'))

			print(data_string, readable)

			fromUser = readable['fromUser']
			toUser   =	readable['toUser']
			assetId = readable['assetId']
			self.send_response(200)
			self.end_headers()
		
			tx, result = bigchainwrapper.transferAsset(fromUser, toUser, assetId)

			self.wfile.write(bytes('{"fromUser":"' + fromUser + '","tx":'+  str(tx) + '}',"utf8"))
			return

try:
	#Create a web server and define the handler to manage the
	#incoming request
	bigchainwrapper.addUser('a')
	bigchainwrapper.addUser('b')
	bigchainwrapper.addUser('c')

	server = HTTPServer(('', PORT_NUMBER), myHandler)
	print ('Started httpserver on port ' , PORT_NUMBER)

	#Wait forever for incoming htto requests
	server.serve_forever()

except KeyboardInterrupt:
	print ('^C received, shutting down the web server')
	server.socket.close()

