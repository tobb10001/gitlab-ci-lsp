import os

from lsprotocol.types import ClientCapabilities, Range
from lsprotocol.types import DefinitionParams
from lsprotocol.types import InitializeParams
from lsprotocol.types import Position
from lsprotocol.types import TextDocumentIdentifier
from lsprotocol.types import Location

import pytest
import pytest_lsp
from pytest_lsp import ClientServerConfig
from pytest_lsp import LanguageClient

SAMPLES_DIR = os.path.join(os.path.dirname(__file__), "..", "samples", "goto_definition")

@pytest_lsp.fixture(
    config=ClientServerConfig(server_command=["go", "run", "main.go"]),
)
async def client(lsp_client: LanguageClient):
    params = InitializeParams(capabilities=ClientCapabilities())
    await lsp_client.initialize_session(params)

    yield

    await lsp_client.shutdown_session()

@pytest.mark.asyncio
async def test_completions(client: LanguageClient):
    dir = os.path.join(SAMPLES_DIR, "local_includes")
    result = await client.text_document_definition_async(
        params=DefinitionParams(
            position=Position(line=1, character=13),
            text_document=TextDocumentIdentifier(uri="file://" + os.path.normpath(os.path.join(dir, ".gitlab-ci.yml")))
        )
    )
    assert isinstance(result, Location)
    assert result.uri == "file://" + os.path.normpath(os.path.join(dir, "directory", ".gitlab-ci.yml"))
    assert result.range == Range(start=Position(line=0, character=0), end=Position(line=0, character=0))
