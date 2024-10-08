CREATE TABLE term (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    term TEXT NOT NULL UNIQUE
);

CREATE TABLE document (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL UNIQUE,
    parser TEXT,
    total_terms INTEGER
);
CREATE TABLE docterm (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL,
    term_id INTEGER NOT NULL,
    term_count INTEGER DEFAULT 1,
    FOREIGN KEY (document_id) REFERENCES document(id) ON DELETE CASCADE,
    FOREIGN KEY (term_id) REFERENCES term(id) ON DELETE CASCADE
);
CREATE TABLE docterm_position (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    docterm_id INTEGER,
    position INTEGER,
    FOREIGN KEY (docterm_id) REFERENCES docterm(id) ON DELETE CASCADE
);


CREATE INDEX docterm_document_id on docterm(document_id);
CREATE INDEX docterm_term_id on docterm(term_id);