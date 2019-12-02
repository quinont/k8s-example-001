import sys
import random
from flask import Flask, request
from flask_restful import reqparse, abort, Api, Resource

app = Flask(__name__)
api = Api(app)

TODOS = {
    'msj1': {'word': 'Buenos D-v2'},
    'msj2': {'word': 'Buas Nes-v2'},
    'msj3': {'word': 'Como est-v2'},
}

def abort_if_todo_doesnt_exist(todo_id):
    if todo_id not in TODOS:
        abort(404, message="Todo {} doesn't exist".format(todo_id))

parser = reqparse.RequestParser()
parser.add_argument('task')

# Todo
#   show a single todo item and lets you delete them
class Todo(Resource):
    def get(self, todo_id):
        if todo_id is None:
           todo_id = TODOS[random.randrange(0, len(TODOS))]  
        sys.stdout.write(str(request.headers))
        sys.stdout.flush()
        abort_if_todo_doesnt_exist(todo_id)
        return TODOS[todo_id]

    def delete(self, todo_id):
        abort_if_todo_doesnt_exist(todo_id)
        del TODOS[todo_id]
        return '', 204

    def put(self, todo_id):
        args = parser.parse_args()
        task = {'task': args['task']}
        TODOS[todo_id] = task
        return task, 201


class Todorand(Resource):
    def get(self):
        todo_id = random.choice(list(TODOS.keys()))
        sys.stdout.write(str(request.headers))
        sys.stdout.flush()
        abort_if_todo_doesnt_exist(todo_id)
        return TODOS[todo_id]


# TodoList
#   shows a list of all todos, and lets you POST to add new tasks
class TodoList(Resource):
    def get(self):
        return TODOS

    def post(self):
        args = parser.parse_args()
        todo_id = 'todo%d' % (len(TODOS) + 1)
        TODOS[todo_id] = {'task': args['task']}
        return TODOS[todo_id], 201

##
## Actually setup the Api resource routing here
##
api.add_resource(TodoList, '/todos')
api.add_resource(Todorand, '/todosrand')
api.add_resource(Todo, '/todos/<string:todo_id>')


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True)
