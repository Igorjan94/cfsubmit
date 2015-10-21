from urls import SERVER_PORT, SERVER_HOST
from valahd_srv.valahd_srv import app

if __name__ == '__main__':
    # HOST = environ.get('SERVER_HOST', 'localhost')
    # try:
    #    PORT = int(environ.get('SERVER_PORT', '5555'))
    # except ValueError:
    #    PORT = 5555
    app.run(SERVER_HOST, SERVER_PORT)
