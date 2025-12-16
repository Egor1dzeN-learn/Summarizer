from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class SummarizeRequest(_message.Message):
    __slots__ = ("text", "prompt")
    TEXT_FIELD_NUMBER: _ClassVar[int]
    PROMPT_FIELD_NUMBER: _ClassVar[int]
    text: str
    prompt: str
    def __init__(self, text: _Optional[str] = ..., prompt: _Optional[str] = ...) -> None: ...

class SummarizeReply(_message.Message):
    __slots__ = ("text",)
    TEXT_FIELD_NUMBER: _ClassVar[int]
    text: str
    def __init__(self, text: _Optional[str] = ...) -> None: ...
