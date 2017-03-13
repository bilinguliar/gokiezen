// Package implements web application with voting by SMS functionality.
//
// DISCLAIMER: Please do not use to elect presidents as there are no guaranties that votes will not be lost due to errors or storage specifics.
//
// But it was built with intention to be quick so there is an option to elect bad presidents really often.
// Note that you will need to pay for outgoing SMS messages and bad decisions.
//
// Current implementation uses Redis as a storage and MessageBird.com as a messaging provider.
// Outbound messages are sent with limited rate of 1 SMS per second. This is a limitation of current provider.
// It utilizes few MessageBird features: receiving SMS, sending SMS and MSISDN lookup.
//
// In order to start Voting you need to add Candidates first. Each candidate can receive votes via short message service.
// You need to send SMS with candidate’s name on a virtual mobile number (VMN) that can be ordered here: https://dashboard.messagebird.com/app/en/numbers
//
// Application exposes an endpoint that will receive POST request with incoming message details.
// Next it will use MessageBird Lookup API call to determine country associated with this MSISDN.
// Score will be updates for candidate, country counter will also be incremented.
package main
