BEGIN;

DROP TABLE IF EXISTS "orders";
CREATE TABLE "public"."orders" (
    "user_id" integer NOT NULL,
    "status" character varying(20) DEFAULT 'NEW' NOT NULL,
    "accrual" double precision DEFAULT '0.0' NOT NULL,
    "created_at" timestamptz DEFAULT now() NOT NULL,
    "id" character(16) NOT NULL,
    CONSTRAINT "orders_id" PRIMARY KEY ("id"),
    CONSTRAINT "orders_status_check" CHECK (((status)::text = ANY ((ARRAY['NEW'::character varying, 'PROCESSING'::character varying, 'INVALID'::character varying, 'PROCESSED'::character varying, 'REGISTERED'::character varying])::text[])))
) WITH (oids = false);


DROP TABLE IF EXISTS "transactions";
DROP SEQUENCE IF EXISTS transactions_id_seq;
CREATE SEQUENCE transactions_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."transactions" (
    "id" integer DEFAULT nextval('transactions_id_seq') NOT NULL,
    "user_id" integer NOT NULL,
    "order_id" character(16) NOT NULL,
    "type_transaction" character varying(15) DEFAULT 'ACCRUAL' NOT NULL,
    "amount" double precision DEFAULT '0.0' NOT NULL,
    "created_at" timestamp DEFAULT now() NOT NULL,
    CONSTRAINT "transactions_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "unq_transaction_unique_columns" UNIQUE ("user_id", "order_id", "type_transaction"),
    CONSTRAINT "transactions_type_transaction_check" CHECK (((type_transaction)::text = ANY ((ARRAY['WITHDRAWAL'::character varying, 'ACCRUAL'::character varying])::text[])))
) WITH (oids = false);


DROP TABLE IF EXISTS "users";
DROP SEQUENCE IF EXISTS users_id_seq;
CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."users" (
    "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "login" character varying(20) NOT NULL,
    "password" character varying(100) NOT NULL,
    "role" character varying(10) DEFAULT 'user' NOT NULL,
    "balance" double precision DEFAULT '0.0' NOT NULL,
    "created_at" timestamp DEFAULT now() NOT NULL,
    CONSTRAINT "users_login_key" UNIQUE ("login"),
    CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "users_role_check" CHECK (((role)::text = ANY ((ARRAY['user'::character varying, 'admin'::character varying])::text[])))
) WITH (oids = false);


ALTER TABLE ONLY "public"."orders" ADD CONSTRAINT "fk_orders_users" FOREIGN KEY (user_id) REFERENCES users(id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."transactions" ADD CONSTRAINT "fk_user_transactions" FOREIGN KEY (user_id) REFERENCES users(id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."transactions" ADD CONSTRAINT "transactions_order_id_fkey" FOREIGN KEY (order_id) REFERENCES orders(id) NOT DEFERRABLE;

COMMIT;