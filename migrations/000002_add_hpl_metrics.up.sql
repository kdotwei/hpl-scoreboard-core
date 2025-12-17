ALTER TABLE "scores" ADD COLUMN "linux_username" varchar NOT NULL DEFAULT '';
ALTER TABLE "scores" ADD COLUMN "n" int NOT NULL DEFAULT 0;
ALTER TABLE "scores" ADD COLUMN "nb" int NOT NULL DEFAULT 0;
ALTER TABLE "scores" ADD COLUMN "p" int NOT NULL DEFAULT 0;
ALTER TABLE "scores" ADD COLUMN "q" int NOT NULL DEFAULT 0;
ALTER TABLE "scores" ADD COLUMN "execution_time" double precision NOT NULL DEFAULT 0.0;

-- Update existing gflops column to be consistent with new structure
-- (keeping it as is for backward compatibility)

-- Add indexes for performance on common query fields
CREATE INDEX ON "scores" ("execution_time");
CREATE INDEX ON "scores" ("n");
CREATE INDEX ON "scores" ("linux_username");