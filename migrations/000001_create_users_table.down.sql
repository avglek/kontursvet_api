-- Indexes
DROP INDEX IF EXISTS  idx_leads_created_at;
DROP INDEX IF EXISTS  idx_portfolio_cases_name ;
DROP INDEX IF EXISTS  idx_portfolio_photos_case_id;
DROP INDEX IF EXISTS  idx_portfolio_case_works_case_id;
DROP INDEX IF EXISTS  idx_portfolio_case_meta_case_id;

-- Leads table
DROP TABLE leads IF EXIST;

-- Portfolio cases
DROP TABLE portfolio_cases IF EXIT;

-- Portfolio case works (many-to-one)
DROP TABLE portfolio_case_works IF EXTIT;

-- Portfolio case meta (many-to-one)
DROP TABLE portfolio_case_meta IF EIXT;

-- Portfolio photos
DROP TABLE  portfolio_photos IF exit;
