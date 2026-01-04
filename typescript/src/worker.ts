import { NativeConnection, Worker } from '@temporalio/worker';
import { loadClientConnectConfig } from '@temporalio/envconfig';
import * as activities from './activities';

async function run() {
  console.info(`⚙️ Using TEMPORAL_PROFILE: '${process.env.TEMPORAL_PROFILE}'`);
  const config = loadClientConnectConfig();
  const connection = await NativeConnection.connect(config.connectionOptions);
  console.info(`✅ Client connected to '${config.connectionOptions.address}' in namespace '${config.namespace}'`);

  try {
    const worker = await Worker.create({
      connection,
      namespace: config.namespace,
      taskQueue: 'trip-task-queue',
      workflowsPath: require.resolve('./workflows'),
      activities,
    });

    await worker.run();
  } finally {
    await connection.close();
  }
}

run().catch((err) => {
  console.error(err);
  process.exit(1);
});
