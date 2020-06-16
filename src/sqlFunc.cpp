/*
 * Filename:	sqlFunc.cpp
 * Author:		Michael Nesbitt
 *
 * Description:
 *
 *	Includes and prototypes for signalHandling.cpp
 */

#include "sqlFunc.h"

static int callback(void *NotUsed, int argc, char **argv, char **azColName)
{
   int i;
   for(i = 0; i<argc; i++)
   {
      writeLog(string(azColName[i]) + " = " + (argv[i] ? argv[i] : "NULL"));
   }
   writeLog("\n", ((long long)NotUsed && false));
   return 0;
}

static int loadIDsCallback(void *NotUsed, int numCols, char** rowData, char **azColName)
{
   serverIDs.push_back(stoi(rowData[0]));
   writeLog(string("Found server: ") + rowData[1] + ", ID: " + rowData[0], ((long long)NotUsed || numCols || azColName || true));
   return 0;
}

int createDB(string dbFile)
{
    sqlite3* db;
    char* zErrMsg = 0;
    int rc;
    string sql;

    rc = sqlite3_open(dbFile.c_str(), &db);

    if( rc )
    {
        writeLog(string("Can't open database: ") + sqlite3_errmsg(db));
        return 1;
    }
    else
    {
        writeLog("Created database successfully");
    }

    /* Create SQL statement */
    sql = "CREATE TABLE Servers( \
        ID INT PRIMARY KEY     NOT NULL, \
        NAME           TEXT    NOT NULL, \
        DIRECTORY      TEXT    NOT NULL, \
        JARFILE        TEXT    NOT NULL \
        );";

    /* Execute SQL statement */
    rc = sqlite3_exec(db, sql.c_str(), callback, 0, &zErrMsg);
    
    if( rc != SQLITE_OK )
    {
        writeLog(string("SQL error: ") + zErrMsg);
        sqlite3_free(zErrMsg);
    }
    else
    {
        writeLog("Table created successfully");
    }

    sqlite3_close(db);

    return 0;
}

int addServerDB(int serverNum, string serverName, string serverDir, string serverJar)
{
    sqlite3* db;
    int rc;

    rc = sqlite3_open(configOptions[DATABASE_FILE].c_str(), &db);

    if( rc )
    {
        writeLog(string("Can't open database: ") + sqlite3_errmsg(db));
        return 1;
    }
    else
    {
        writeLog("Opened database successfully");
    }

    string sql = "INSERT INTO Servers (ID, NAME, DIRECTORY, JARFILE) "
                 "VALUES (" + to_string(serverNum) + ", \"" + serverName + "\", \"" + serverDir + "\", \"" + serverJar + "\");";
    char* zErrMsg = 0;

    rc = sqlite3_exec(db, sql.c_str(), NULL, 0, &zErrMsg);
    if (rc != SQLITE_OK)
    {
        writeLog("Error inserting into database: " + string(zErrMsg));
        sqlite3_free(zErrMsg);
        return 1;
    } 
    else
    {
        writeLog("Record created successfully!");
    }

    sqlite3_close(db);

    return 0;
}

int loadIDs()
{
    sqlite3* db;
    int rc;

    rc = sqlite3_open(configOptions[DATABASE_FILE].c_str(), &db);

    if( rc )
    {
        writeLog(string("Can't open database: ") + sqlite3_errmsg(db));
        return 1;
    }
    else
    {
        writeLog("Opened database successfully");
    }

    string sql = "SELECT * FROM Servers;";
    char* zErrMsg = 0;

    rc = sqlite3_exec(db, sql.c_str(), loadIDsCallback, 0, &zErrMsg);
    if (rc != SQLITE_OK)
    {
        writeLog("Error querying database: " + string(zErrMsg));
        sqlite3_free(zErrMsg);
        return 1;
    } 
    else
    {
        writeLog("Records loaded successfully!");
    }

    sqlite3_close(db);

    return 0;
}
