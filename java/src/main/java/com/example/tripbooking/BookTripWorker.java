package com.example.tripbooking;

import com.example.tripbooking.activities.TripActivitiesImpl;
import com.example.tripbooking.workflows.BookWorkflowImpl;
import io.temporal.client.WorkflowClient;
import io.temporal.client.WorkflowClientOptions;
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
        logger.info("Starting Temporal worker...");

        // Create service stubs
        WorkflowServiceStubs serviceStubs = WorkflowServiceStubs.newServiceStubs(
            WorkflowServiceStubsOptions.newBuilder().setTarget("localhost:7233").build()
        );

        // Create workflow client
        WorkflowClient client = WorkflowClient.newInstance(
            serviceStubs,
            WorkflowClientOptions.newBuilder().setNamespace("default").build()
        );

        // Create worker factory
        WorkerFactory factory = WorkerFactory.newInstance(client);

        // Create worker
        Worker worker = factory.newWorker(TASK_QUEUE);

        // Register workflow and activities
        worker.registerWorkflowImplementationTypes(BookWorkflowImpl.class);
        worker.registerActivitiesImplementations(new TripActivitiesImpl());

        // Start the worker
        factory.start();
    }
}
