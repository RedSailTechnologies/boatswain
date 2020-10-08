// package: 
// file: kraken/service.proto

var kraken_service_pb = require("../kraken/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var Kraken = (function () {
  function Kraken() {}
  Kraken.serviceName = "Kraken";
  return Kraken;
}());

Kraken.Clusters = {
  methodName: "Clusters",
  service: Kraken,
  requestStream: false,
  responseStream: false,
  requestType: kraken_service_pb.ClusterRequest,
  responseType: kraken_service_pb.ClusterResponse
};

Kraken.Deployments = {
  methodName: "Deployments",
  service: Kraken,
  requestStream: false,
  responseStream: false,
  requestType: kraken_service_pb.DeploymentRequest,
  responseType: kraken_service_pb.DeploymentResponse
};

exports.Kraken = Kraken;

function KrakenClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

KrakenClient.prototype.clusters = function clusters(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Kraken.Clusters, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

KrakenClient.prototype.deployments = function deployments(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Kraken.Deployments, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.KrakenClient = KrakenClient;

