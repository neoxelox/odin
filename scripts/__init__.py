import superinvoke

from . import tasks
from .tools import Tools

root = superinvoke.init(Tools)

root.add_task(tasks.start)
root.add_task(tasks.clean)
root.add_task(tasks.prune)
root.add_task(tasks.shell)
root.add_task(tasks.cli)
