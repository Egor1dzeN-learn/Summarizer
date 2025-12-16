from dependency_injector import containers, providers

from .services import ServiceContainer


class RootContainer(containers.DeclarativeContainer):
  config = providers.Configuration()

  services = providers.Container(ServiceContainer)
