DROP TABLE IF EXISTS Requests;
CREATE TABLE Requests
(
    vt_request_id          TEXT PRIMARY KEY,
    float_timeoff_id       INTEGER NOT NULL,
    created                INTEGER NOT NULL
);
