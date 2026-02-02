from collections.abc import AsyncGenerator

import pytest_asyncio
from temporalio.client import Client
from temporalio.testing import WorkflowEnvironment


@pytest_asyncio.fixture(scope="session")
async def env() -> AsyncGenerator[WorkflowEnvironment, None]:
    env = await WorkflowEnvironment.start_time_skipping()
    yield env
    await env.shutdown()


@pytest_asyncio.fixture
async def client(env: WorkflowEnvironment) -> Client:
    return env.client
