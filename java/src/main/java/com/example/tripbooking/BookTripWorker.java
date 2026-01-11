package com.example.tripbooking;

import com.example.tripbooking.activities.TripActivitiesImpl;
import com.example.tripbooking.workflows.BookWorkflowImpl;
import io.temporal.client.WorkflowClient;
import io.temporal.client.WorkflowClientOptions;
import io.temporal.envconfig.ClientConfigProfile;
import io.temporal.serviceclient.WorkflowServiceStubs;
import io.temporal.serviceclient.WorkflowServiceStubsOptions;
import io.temporal.worker.Worker;
import io.temporal.worker.WorkerFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class BookTripWorker {

    private static final Logger logger = LoggerFactory.getLogger(BookTripWorker.class);
    private static final String TASK_QUEUE = "trip-task-queue";

    public static void main(String[] args) {
        try {
            logger.info("⚙️ Using TEMPORAL_PROFILE: '{}'", System.getenv("TEMPORAL_PROFILE"));
            ClientConfigProfile profile = ClientConfigProfile.load();
            WorkflowServiceStubsOptions serviceStubsOptions = profile.toWorkflowServiceStubsOptions();
            WorkflowClientOptions clientOptions = profile.toWorkflowClientOptions();
            WorkflowServiceStubs serviceStubs = WorkflowServiceStubs.newServiceStubs(serviceStubsOptions);
            WorkflowClient client = WorkflowClient.newInstance(serviceStubs, clientOptions);
            WorkerFactory factory = WorkerFactory.newInstance(client);

            Worker worker = factory.newWorker(TASK_QUEUE);
            worker.registerWorkflowImplementationTypes(BookWorkflowImpl.class);
            worker.registerActivitiesImplementations(new TripActivitiesImpl());
            factory.start();
            logger.info("✅ Client connected to '{}' in namespace '{}'",
                    serviceStubsOptions.getTarget(), clientOptions.getNamespace());
            logger.info("Worker started and listening on task queue: '{}'", TASK_QUEUE);
        } catch (Exception e) {
            logger.error("Failed to start Temporal worker: {}", e.getMessage(), e);
            System.exit(1);
        }
    }
}
