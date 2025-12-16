from dependency_injector import containers, providers

from worker.impl.services import TransformerSummarizerService


class ServiceContainer(containers.DeclarativeContainer):
  summarizer = providers.Singleton(TransformerSummarizerService)
