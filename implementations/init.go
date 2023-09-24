// Package implementations serves as the core of the database interaction framework.
// It provides implementations of various interfaces and structures necessary for
// seamless communication with different database systems.
//
// Key Components:
// - dbCommon: This package defines a common structure 'DBCommon' to share
//   common functionality among various database implementations.
//
// - SqlBase: The 'SqlBase' struct contains implementations of SQL methods that are
//   shared by different SQL-based databases. It serves as a foundation for
//   database-specific implementations.
//
// - MySql: The 'MySql' struct extends 'SqlBase' and implements functionality
//   specific to MySQL databases, including query construction and execution.
//
// - MongoConnection: The 'MongoConnection' struct implements functionality for
//   interacting with MongoDB databases, such as connection management, query
//   building, and aggregation pipelines.
//
// Purpose:
// This package offers a flexible and modular approach to working with databases.
// It separates common and database-specific functionality, allowing developers
// to build applications that can seamlessly switch between different database
// systems by changing the database type without major code modifications.
//
//
// Note:
// Developers can build upon this package to create database-agnostic applications
// that can seamlessly adapt to different database systems.
package implementations
