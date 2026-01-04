from temporalio.client import Client
from temporalio.envconfig import ClientConfigProfile


async def get_client():
    # Load the "default" profile from default locations and environment variables.
    # Environment variables take precedence over configuration file settings.
    default_profile = ClientConfigProfile.load()
    connect_config = default_profile.to_client_connect_config()

    # Connect to the client using the loaded configuration.
    client = await Client.connect(**connect_config)
    print(f"✅ Client connected to {client.service_client.config.target_host} in namespace '{client.namespace}'")
    return client
