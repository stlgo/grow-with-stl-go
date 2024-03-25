# WebService Clients

Let's take a stroll down random access memory lane to see where we came from and where we're going.

## First there was the Common Gateway Interface (CGI)

First proposed in 1993 CGI became [RFC 3875](https://datatracker.ietf.org/doc/html/rfc3875) in 1997.  Web Servers like [Apache HTTPD](https://httpd.apache.org/) allowed for C and other language programs to be executed when called via a web request allowing for dynamic information to be sent as a response.

## Then came SOAP (Simple Object Access Protocol)

First appearing as [XML-RPC](https://en.wikipedia.org/wiki/XML-RPC) in 1998 and later becoming [RFC 4227](https://datatracker.ietf.org/doc/html/rfc4227) which we know today as SOAP.  SOAP can be thought of as XML over HTTP, it wes conceived as a machine to machine protocol and allowed for the first wide spread use of what we know of as web services.  In order to use SOAP a client would first need to get the [WSDL (Web Service Description Language)](https://www.w3schools.com/xml/xml_wsdl.asp) that would define the endpoints, inputs and responses from a given web service.  This would need to be compiled into a client which could then flex the web service.  While this is an improvement over CGI it requires a strict interface agreement and is harder to interact with the web service than what protocols that superseded it.

## REST because SOAP wasn't as good as it should have been

REST (representational state transfer) first appeared in Roy Fielding's 2000 PhD dissertation later becoming [RFC 6690](https://datatracker.ietf.org/doc/html/rfc6690).  It was conceived of in response to the short comings of SOAP.  It allows for the use of various inputs and outputs over the standard [HTTP methods](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods) GET, PUT, POST, DELETE being the most common and HEAD, CONNECT, OPTIONS, TRACE, PATCH being less common methods.  Often used in conjunction with [JSON (JavaScript Object Notation)](https://www.json.org/json-en.html) it forms the basis of most of the web services that exist today.

## WebSocket

WebSockets first emerged around 2010 and later became [RFC 6455](https://datatracker.ietf.org/doc/html/rfc6455).  WebSocket handshake uses the HTTP Upgrade header to change from the HTTP protocol to the WebSocket protocol.  WebSockets are fully duplexed meaning it can both send and receive messages at the same time.  With REST you must poll to get updates, with WebSockets you can get information when it is available and while REST is a great protocol, WebSockets has advantages and while not as ubiquitous as REST it certainly should be used whenever real time communication and processing is a must in an API setting.

## Grow With STL-Go Examples

Because the Grow With STL-Go application implements both a REST web service and WebSocket we can flex the most common APIs today and one that enables real time interaction.

- [REST](REST/README.md) client example
- [WebSocket](WebSocket/README.md) client example
