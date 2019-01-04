"use strict";
var ids = {
  serverAddress: "inputAddress",
  outputFormatted: "spanOutputFormatted",
  outputRaw: "spanOutputRaw",
  key: "inputKey",
  value: "inputValue",
};

function processRequest(input) {
  document.getElementById(ids.outputRaw).innerHTML = input;
  var inputParsed = JSON.parse( input);
  var resultHTML = "";
  if (inputParsed.error !== undefined && inputParsed.error !== null) {
    resultHTML += `<b>Error:</b> <b style='color:red'>${inputParsed.error}</b><br>`;
  }  
  if (inputParsed.comments !== undefined && inputParsed.comments !== null) {
    resultHTML += `<b>Comments:</b> ${inputParsed.comments}<br>`;
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
  };
  request.open("GET", `/request?json=${encodeURIComponent(JSON.stringify(jsonRequest))}`, true);
  request.onload = function() {
    processRequest(request.responseText);
  }
  ////////////////////////////////////////////
  request.send(null);
} 
