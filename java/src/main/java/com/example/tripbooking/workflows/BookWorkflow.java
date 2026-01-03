package com.example.tripbooking.workflows;

import com.example.tripbooking.shared.BookTripInput;

import io.temporal.workflow.WorkflowInterface;
import io.temporal.workflow.WorkflowMethod;

@WorkflowInterface
public interface BookWorkflow {
    @WorkflowMethod
    String bookTrip(BookTripInput input);
}
