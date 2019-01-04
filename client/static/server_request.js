"use strict";
var ids = {
  serverAddress: "inputAddress",
  outputFormatted: "spanOutputFormatted",
  outputRaw: "spanOutputRaw",
  key: "inputKey",
  value: "inputValue",
  peer: "inputPeer",
};

function processRequest(input) {
  document.getElementById(ids.outputRaw).innerHTML = input;
  var inputParsed = JSON.parse( input);
  var resultHTML = "";

  if (inputParsed.error !== undefined && inputParsed.error !== null) {
    resultHTML += `<b>Client error:</b> <b style='color:red'>${inputParsed.error}</b><br>`;
  }  
  if (inputParsed.result !== undefined) {
    if (inputParsed.result.error !== undefined) {
      resultHTML += `<b>Server error: </b> <b style='color:red'>${inputParsed.result.error}</b><br>`
    }
    var data = inputParsed.result.data; 
    if (data !== undefined) {
      resultHTML += `<b>Data:</b> `;
      resultHTML += `<table border="1"><tr><th>key</th><th>value</th><th>Version</th></tr>`;
      for (var label in data) {
        resultHTML += `<tr><td>${label}</td><td>${data[label].value}</td><td>${data[label].version}</td></tr>`;
      }
      resultHTML += "</table><br>";
    }
    var serverAddress = inputParsed.result.serverAddress;
    if (serverAddress !== undefined) {
      resultHTML += `<b>Server address:</b> ${serverAddress}<br>`;
    }
    if (inputParsed.result.comments !== undefined) {
      resultHTML += `<b>Server comments:</b> ${inputParsed.result.comments}<br>`;
    }
  }

  if (inputParsed.comments !== undefined && inputParsed.comments !== null) {
    resultHTML += `<b>Client comments:</b> ${inputParsed.comments}<br>`;
  }
  document.getElementById(ids.outputFormatted).innerHTML = resultHTML;
}

function getRequest(inputCommand) {
  var request = new XMLHttpRequest();
  var jsonRequest = {
    command: inputCommand,
    serverAddress: document.getElementById(ids.serverAddress).value,
    key: document.getElementById(ids.key).value,
    value: document.getElementById(ids.value).value,
    peer: document.getElementById(ids.peer).value,
    version: (new Date()).getTime(),
  };
  request.open("GET", `/request?json=${encodeURIComponent(JSON.stringify(jsonRequest))}`, true);
  request.onload = function() {
    processRequest(request.responseText);
  }
  ////////////////////////////////////////////
  request.send(null);
} 
