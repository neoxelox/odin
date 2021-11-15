from superinvoke import task

Composer = "docker-compose --env-file .env.dev --file docker-compose.dev.yaml"


@task()
def start(context):
    """Start the project."""
    context.run(f"{Composer} up --build")

@task()
def clean(context):
    """Remove all containers, volumes and networks of the project."""
    context.run(f"{Composer} down --volumes")

@task()
def prune(context):
    """Remove all containers, volumes, networks and images of the system."""
    context.run("docker system prune -a --volumes")

@task()
def shell(context):
    """Enter the development shell."""
    context.run(f"{Composer} exec --workdir /workspace dev bash", pty=True)

@task(variadic=True)
def cli(context, args):
    """Execute Odin's commands."""
    context.run(f"{Composer} exec -T odin-worker cli {args}")
