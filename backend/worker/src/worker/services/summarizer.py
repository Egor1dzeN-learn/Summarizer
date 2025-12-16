from abc import ABC, abstractmethod


class SummarizerService(ABC):
  @abstractmethod
  def summarize(self, text: str, prompt: str) -> str:
    pass
