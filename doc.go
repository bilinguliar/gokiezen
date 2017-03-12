// Package implements web application with voting by SMS functionality.
//
// DISCLAIMER: Please do not use to elect presidents as there are no guaranties that votes will not be lost due to errors and storage specifics.
//
// But it was with intention to be quick so there is an option to elect bad presidents really often.
// Note that you will need to pay for outgoing SMS messages and bad descisions.
//
// Current implementation uses Redis as a storage and MessageBird.com as a messaging provider.
// It utilizes few MessageBird features: reciving SMS, sending SMS and MSISDN lookup.
//
// In order to start Voting you need to add Candidates first. Each candidate can recieve votes via short message service.
// You need to send SMS with candidate name on a virtual mobile number that can be ordered here: https://dashboard.messagebird.com/app/en/numbers
//
// Application exposes an endpoint that will recieve POST request with incomming message details.
// Next it will use MessageBird Lookup API call to determine country associated with this MSISDN.
// Score will be updates for candidate, country counter will also be incremented.
package main
