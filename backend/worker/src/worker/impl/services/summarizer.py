from transformers import T5ForConditionalGeneration, T5Tokenizer
from worker.services import SummarizerService
import os
token = os.getenv("HF_TOKEN")

class TransformerSummarizerService(SummarizerService):
  def __init__(self):
    path_to_model = "oOundefinedOo/QA-system-T5_RUS"

    self._model = T5ForConditionalGeneration.from_pretrained(
      path_to_model,
      use_safetensors=True,
      token=token,
      cache_dir="C:/hf_cache",
    )
    self._tokenizer = T5Tokenizer.from_pretrained(
      path_to_model,
      token=token,
      cache_dir="C:/hf_cache",
    )

  def summarize(self, text: str, prompt: str) -> str:
    input_text = f"question: {prompt} context: {text}"
    input_ids = self._tokenizer.encode(input_text, return_tensors="pt")
    outputs = self._model.generate(input_ids, max_length=20)
    answer = self._tokenizer.decode(outputs[0], skip_special_tokens=True)
    return answer
