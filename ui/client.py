import os

from temporalio.client import Client
from temporalio.envconfig import ClientConfig


async def get_client():
    print(f"⚙️ Using TEMPORAL_PROFILE: '{os.environ.get('TEMPORAL_PROFILE')}'")
    connect_config = ClientConfig.load_client_connect_config()
    client = await Client.connect(**connect_config)
    print(f"✅ Client connected to '{client.service_client.config.target_host}' in namespace '{client.namespace}'")
    return client
