package org.example.service.services;

public class Result<T, U>  {
    private T result; // result not null if success
    private U returnStatus;
    public Result() {
        result = null;
        returnStatus = null;
    }

    public Result(U returnStatus) {
        this.returnStatus = returnStatus;
    }

    public T getResult() {
        return result;
    }

    public void setResult(T result) {
        this.result = result;
    }

    public U getReturnStatus() {
        return returnStatus;
    }

    public void setReturnStatus(U returnStatus) {
        this.returnStatus = returnStatus;
    }
}
