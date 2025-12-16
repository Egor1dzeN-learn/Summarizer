import os

from dotenv import load_dotenv

from worker.di import RootContainer
from worker.di.services import ServiceContainer
from worker.server import serve

if __name__ == "__main__":
  load_dotenv()

  container = ServiceContainer()  # todo: fix & use RootContainer
  container.wire(modules=["worker.server"])

  serve(
    host=os.environ.get("APP_HOST", "localhost"),
    port=int(os.environ.get("APP_PORT", 50051)),
  )
