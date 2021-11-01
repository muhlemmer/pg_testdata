create table unit_tests (
    bool_col	bool	null
);

create table error_tests (
    bool_col    bool    null,
    text_col    text    not null
);

create table regression_tests (
    bool_col_n  bool    null,
    bool_col_nn bool    not null
);