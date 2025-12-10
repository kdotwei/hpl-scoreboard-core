CREATE TABLE "scores" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" varchar NOT NULL, 
  "gflops" double precision NOT NULL,
  "problem_size_n" int NOT NULL,
  "block_size_nb" int NOT NULL,
  "submitted_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "scores" ("gflops");