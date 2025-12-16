from concurrent import futures

import grpc
from dependency_injector.wiring import Provide, inject
from loguru import logger

from worker.di import ServiceContainer
from worker.generated.protos import worker_pb2, worker_pb2_grpc
from worker.services.summarizer import SummarizerService


class NodeWorker(worker_pb2_grpc.NodeWorkerServicer):
  @inject
  def __init__(
    self,
    svc: SummarizerService = Provide[ServiceContainer.summarizer],
  ):
    self._svc = svc

  def Summarize(
    self,
    request: worker_pb2.SummarizeRequest,
    context: grpc.ServicerContext,
  ):
    logger.debug("req=({})", request)
    return worker_pb2.SummarizeReply(
      text=self._svc.summarize(request.text, request.prompt)
    )


def serve(host: str, port: int):
  server = grpc.server(futures.ThreadPoolExecutor(max_workers=3))
  server.add_insecure_port(f"{host}:{port}")

  w = NodeWorker()
  worker_pb2_grpc.add_NodeWorkerServicer_to_server(w, server)

  logger.info("Listening on {}:{}", host, port)
  server.start()
  server.wait_for_termination()
