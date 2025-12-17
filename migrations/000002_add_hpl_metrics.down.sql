DROP INDEX IF EXISTS scores_execution_time_idx;
DROP INDEX IF EXISTS scores_n_idx;  
DROP INDEX IF EXISTS scores_linux_username_idx;

ALTER TABLE "scores" DROP COLUMN IF EXISTS "linux_username";
ALTER TABLE "scores" DROP COLUMN IF EXISTS "n";
ALTER TABLE "scores" DROP COLUMN IF EXISTS "nb";
ALTER TABLE "scores" DROP COLUMN IF EXISTS "p";
ALTER TABLE "scores" DROP COLUMN IF EXISTS "q";
ALTER TABLE "scores" DROP COLUMN IF EXISTS "execution_time";