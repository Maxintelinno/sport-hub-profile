ระบบต้องมี 3 เรื่องแยกกันชัด ๆ
1. ข้อมูลบัญชีรับเงินของ owner
2. ยอดที่ owner ควรได้รับ
3. ประวัติการโอนเงินจริง

โครง flow ที่แนะนำ
Flow เงิน

ลูกค้าจ่ายเงิน
→ payment สำเร็จ
→ ระบบสร้าง settlement ให้ owner
→ owner มียอดค้างรับ
→ ถึงรอบโอนเงิน / owner กดขอถอน
→ ระบบสร้าง payout
→ โอนเงินจริง
→ บันทึกผลการโอน

1) ตารางบัญชีรับเงินของ owner
ผมแนะนำตาราง owner_bank_accounts

CREATE TABLE public.owner_bank_accounts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,

    bank_code varchar(20) NOT NULL,
    bank_name varchar(100) NOT NULL,
    account_name varchar(255) NOT NULL,
    account_number varchar(50) NOT NULL,

    promptpay_type varchar(20) NULL, 
    -- phone / national_id / ewallet / null
    promptpay_value varchar(100) NULL,

    is_default boolean NOT NULL DEFAULT false,
    is_verified boolean NOT NULL DEFAULT false,
    verification_status varchar(30) NOT NULL DEFAULT 'pending',
    -- pending / verified / rejected

    status varchar(20) NOT NULL DEFAULT 'active',
    -- active / inactive

    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT owner_bank_accounts_pkey PRIMARY KEY (id),
    CONSTRAINT owner_bank_accounts_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES public.users(id)
);

ฟิลด์สำคัญ
user_id = เจ้าของสนามคนไหน
bank_code = เช่น KBANK, SCB, BBL
account_name = ชื่อบัญชี
account_number = เลขบัญชี
is_default = บัญชีหลัก
is_verified = ผ่านการตรวจสอบหรือยัง
 
#######################
ควรเก็บ PromptPay ไหม
#######################

ควรครับ ถ้าอนาคตคุณอยากโอนผ่าน PromptPay

เช่น:

เบอร์โทร
เลขบัตร
wallet id

แต่ถ้าตอนนี้ยังไม่ใช้ จะปล่อย NULL ไว้ก่อนก็ได้

เรื่อง security สำคัญมาก
เลขบัญชีถือเป็นข้อมูลการเงิน ควรอย่างน้อย:

แสดงใน FE แบบ mask เช่น xxx-x-12345-x
เก็บ encrypted at rest ถ้าทำได้
อย่า log เลขบัญชีเต็มใน application log

--------------------
2) ตาราง settlement

อันนี้คือ “เงินที่ owner ควรได้รับ” จากแต่ละ booking

CREATE TABLE public.owner_settlements (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    booking_id uuid NOT NULL,
    owner_id uuid NOT NULL,

    gross_amount numeric(10,2) NOT NULL,
    platform_fee numeric(10,2) NOT NULL DEFAULT 0,
    discount_amount numeric(10,2) NOT NULL DEFAULT 0,
    net_amount numeric(10,2) NOT NULL,

    status varchar(30) NOT NULL DEFAULT 'pending',
    -- pending / available / processing / paid / hold / reversed

    available_at timestamp NULL,
    paid_at timestamp NULL,

    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT owner_settlements_pkey PRIMARY KEY (id),
    CONSTRAINT owner_settlements_booking_id_fkey
        FOREIGN KEY (booking_id) REFERENCES public.bookings(id),
    CONSTRAINT owner_settlements_owner_id_fkey
        FOREIGN KEY (owner_id) REFERENCES public.users(id),
    CONSTRAINT owner_settlements_booking_unique UNIQUE (booking_id)
);

ทำไมต้องมี settlement

เพราะ payment ที่ลูกค้าจ่าย ≠ เงินสุทธิที่ owner ได้

ตัวอย่าง:
ลูกค้าจ่าย 1,000
ระบบหัก platform fee 100
owner ได้ 900

ดังนั้นต้องมีตารางนี้เพื่อเป็น source of truth

3) ตาราง payout

อันนี้คือ “รอบโอนเงินจริง”

CREATE TABLE public.owner_payouts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    payout_no varchar(30) NOT NULL,
    owner_id uuid NOT NULL,
    bank_account_id uuid NOT NULL,

    total_amount numeric(10,2) NOT NULL,
    transfer_fee numeric(10,2) NOT NULL DEFAULT 0,
    final_amount numeric(10,2) NOT NULL,

    status varchar(30) NOT NULL DEFAULT 'pending',
    -- pending / processing / paid / failed / cancelled

    payout_method varchar(30) NOT NULL DEFAULT 'bank_transfer',
    -- bank_transfer / promptpay

    requested_at timestamp NULL,
    processed_at timestamp NULL,
    paid_at timestamp NULL,
    failed_at timestamp NULL,
    failure_reason text NULL,

    transfer_reference varchar(150) NULL,
    transfer_provider varchar(50) NULL,

    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT owner_payouts_pkey PRIMARY KEY (id),
    CONSTRAINT owner_payouts_payout_no_key UNIQUE (payout_no),
    CONSTRAINT owner_payouts_owner_id_fkey
        FOREIGN KEY (owner_id) REFERENCES public.users(id),
    CONSTRAINT owner_payouts_bank_account_id_fkey
        FOREIGN KEY (bank_account_id) REFERENCES public.owner_bank_accounts(id)
);

4) ตารางเชื่อม settlement กับ payout 

เพราะ 1 รอบ payout อาจรวมหลาย booking

CREATE TABLE public.owner_payout_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    payout_id uuid NOT NULL,
    settlement_id uuid NOT NULL,
    amount numeric(10,2) NOT NULL,

    created_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT owner_payout_items_pkey PRIMARY KEY (id),
    CONSTRAINT owner_payout_items_payout_id_fkey
        FOREIGN KEY (payout_id) REFERENCES public.owner_payouts(id),
    CONSTRAINT owner_payout_items_settlement_id_fkey
        FOREIGN KEY (settlement_id) REFERENCES public.owner_settlements(id),
    CONSTRAINT owner_payout_items_unique UNIQUE (payout_id, settlement_id)
);

Flow การทำงานที่แนะนำ
เมื่อ booking จ่ายสำเร็จ
payment = paid
สร้าง owner_settlements
status = pending หรือ available
เมื่อถึงรอบโอน

เช่นโอนทุกวันศุกร์ หรือ owner กดถอนเอง

หา settlements ที่ available
รวมยอด
สร้าง owner_payouts
สร้าง owner_payout_items
อัปเดต settlements เป็น processing
เมื่อโอนสำเร็จ
owner_payouts.status = paid
owner_settlements.status = paid
set paid_at
ถ้าโอนล้มเหลว
owner_payouts.status = failed
owner_settlements.status กลับเป็น available
เก็บ failure_reason
ควรให้ owner ถอนเอง หรือโอนอัตโนมัติ

มี 2 แบบ

แบบ A: โอนอัตโนมัติ

เช่นทุกสัปดาห์ / ทุกเดือน

ข้อดี:

owner ไม่ต้องกดอะไร
ระบบดูโปร

ข้อเสีย:

backend ซับซ้อนขึ้น
ต้องมี scheduler
แบบ B: owner กด “ขอถอนเงิน”

ข้อดี:

ง่ายกว่า
คุมง่ายช่วงแรก

ข้อเสีย:

owner ต้องมี action เอง
สำหรับตอนนี้

ผมแนะนำ:

เริ่มจาก manual request payout ก่อน
แล้วค่อยไป auto payout ตอนระบบนิ่ง

ควรมี field เพิ่มใน users ไหม

ถ้าจะให้สะดวก อาจเพิ่ม field บอกว่า owner ตั้งบัญชีหรือยัง

ALTER TABLE public.users
ADD COLUMN has_payout_account boolean NOT NULL DEFAULT false;

แต่ อย่าเก็บเลขบัญชีใน users ตรง ๆ
ให้เก็บใน owner_bank_accounts

API ที่ควรมี
จัดการบัญชี owner
GET    /v1/owner/bank-accounts
POST   /v1/owner/bank-accounts
PUT    /v1/owner/bank-accounts/{id}
DELETE /v1/owner/bank-accounts/{id}
POST   /v1/owner/bank-accounts/{id}/set-default
ดูยอดค้างรับ
GET /v1/owner/wallet/summary

response:

{
  "available_balance": 3200,
  "processing_balance": 900,
  "paid_out_total": 15000
}
ขอถอนเงิน
POST /v1/owner/payouts/request

body:

{
  "bank_account_id": "uuid"
}
ดูประวัติการโอน
GET /v1/owner/payouts
GET /v1/owner/payouts/{id}
Validation ที่ควรมี

ก่อนให้ owner ขอถอน:

ต้องมีบัญชี default
บัญชีต้อง verified
ยอด available ต้องมากกว่าขั้นต่ำ เช่น 100 บาท
ต้องไม่มี payout processing อยู่ซ้อนกัน
Verification บัญชีควรทำยังไง

เริ่มง่ายสุด:

owner กรอกชื่อบัญชี + เลขบัญชี
admin ตรวจสอบ manual
set is_verified = true

โตขึ้นค่อยทำ:

bank account name match
penny test transfer
KYC
ถ้าอยากรองรับภาษี/เอกสารในอนาคต

อาจเตรียม field เพิ่มใน owner profile แยกอีกตาราง เช่น:

เลขบัตร/เลขผู้เสียภาษี
ที่อยู่สำหรับออกเอกสาร
ประเภทบุคคล/นิติบุคคล

แต่ตอนนี้ยังไม่จำเป็นถ้ายัง MVP

สรุปโครงสร้างที่ควรมี
ตารางใหม่
owner_bank_accounts
owner_settlements
owner_payouts
owner_payout_items
users

อาจเพิ่มแค่

has_payout_account
สรุปแบบง่ายที่สุด

ถ้า owner คนหนึ่งมีบัญชี 1 บัญชี ก็ยังควรแยกตารางนะ เพราะในอนาคตจะได้:

เปลี่ยนบัญชี
มีหลายบัญชี
ตรวจสอบสถานะบัญชี
ดูประวัติ payout ได้สะอาด

1) SQL migration final ครบชุด
1.1 เพิ่ม flag ใน users

เอาไว้เช็กเร็ว ๆ ว่า owner ตั้งบัญชีรับเงินแล้วหรือยัง
ALTER TABLE public.users
ADD COLUMN has_payout_account boolean NOT NULL DEFAULT false;

1.2 ตารางบัญชีรับเงินของ owner

CREATE TABLE public.owner_bank_accounts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,

    bank_code varchar(20) NOT NULL,
    bank_name varchar(100) NOT NULL,
    account_name varchar(255) NOT NULL,
    account_number varchar(50) NOT NULL,

    promptpay_type varchar(20) NULL,
    -- phone / national_id / ewallet / null

    promptpay_value varchar(100) NULL,

    is_default boolean NOT NULL DEFAULT false,
    is_verified boolean NOT NULL DEFAULT false,
    verification_status varchar(30) NOT NULL DEFAULT 'pending',
    -- pending / verified / rejected

    status varchar(20) NOT NULL DEFAULT 'active',
    -- active / inactive

    note text NULL,

    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT owner_bank_accounts_pkey PRIMARY KEY (id),
    CONSTRAINT owner_bank_accounts_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);

CREATE INDEX idx_owner_bank_accounts_user_id
    ON public.owner_bank_accounts(user_id);

CREATE INDEX idx_owner_bank_accounts_user_default
    ON public.owner_bank_accounts(user_id, is_default);

CREATE INDEX idx_owner_bank_accounts_status
    ON public.owner_bank_accounts(status);

1.3 ตาราง settlement

เก็บยอดที่ owner “ควรได้รับ” ต่อ booking
CREATE TABLE public.owner_settlements (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    booking_id uuid NOT NULL,
    owner_id uuid NOT NULL,

    gross_amount numeric(10,2) NOT NULL,
    platform_fee numeric(10,2) NOT NULL DEFAULT 0,
    discount_amount numeric(10,2) NOT NULL DEFAULT 0,
    net_amount numeric(10,2) NOT NULL,

    status varchar(30) NOT NULL DEFAULT 'pending',
    -- pending / available / processing / paid / hold / reversed

    available_at timestamp NULL,
    paid_at timestamp NULL,

    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT owner_settlements_pkey PRIMARY KEY (id),
    CONSTRAINT owner_settlements_booking_id_fkey
        FOREIGN KEY (booking_id) REFERENCES public.bookings(id) ON DELETE CASCADE,
    CONSTRAINT owner_settlements_owner_id_fkey
        FOREIGN KEY (owner_id) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT owner_settlements_booking_unique UNIQUE (booking_id)
);

CREATE INDEX idx_owner_settlements_owner_id
    ON public.owner_settlements(owner_id);

CREATE INDEX idx_owner_settlements_owner_status
    ON public.owner_settlements(owner_id, status);

CREATE INDEX idx_owner_settlements_available_at
    ON public.owner_settlements(available_at);

1.4 ตาราง payout

เก็บรอบโอนเงินจริง
CREATE TABLE public.owner_payouts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    payout_no varchar(30) NOT NULL,
    owner_id uuid NOT NULL,
    bank_account_id uuid NOT NULL,

    total_amount numeric(10,2) NOT NULL,
    transfer_fee numeric(10,2) NOT NULL DEFAULT 0,
    final_amount numeric(10,2) NOT NULL,

    status varchar(30) NOT NULL DEFAULT 'pending',
    -- pending / processing / paid / failed / cancelled

    payout_method varchar(30) NOT NULL DEFAULT 'bank_transfer',
    -- bank_transfer / promptpay

    requested_at timestamp NULL,
    processed_at timestamp NULL,
    paid_at timestamp NULL,
    failed_at timestamp NULL,
    cancelled_at timestamp NULL,

    failure_reason text NULL,

    transfer_reference varchar(150) NULL,
    transfer_provider varchar(50) NULL,

    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT owner_payouts_pkey PRIMARY KEY (id),
    CONSTRAINT owner_payouts_payout_no_key UNIQUE (payout_no),
    CONSTRAINT owner_payouts_owner_id_fkey
        FOREIGN KEY (owner_id) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT owner_payouts_bank_account_id_fkey
        FOREIGN KEY (bank_account_id) REFERENCES public.owner_bank_accounts(id)
);


CREATE INDEX idx_owner_payouts_owner_id
    ON public.owner_payouts(owner_id);

CREATE INDEX idx_owner_payouts_owner_status
    ON public.owner_payouts(owner_id, status);

CREATE INDEX idx_owner_payouts_requested_at
    ON public.owner_payouts(requested_at);

1.5 ตาราง payout items

เชื่อม settlement หลายรายการเข้ากับ payout เดียว
CREATE TABLE public.owner_payout_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    payout_id uuid NOT NULL,
    settlement_id uuid NOT NULL,
    amount numeric(10,2) NOT NULL,

    created_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT owner_payout_items_pkey PRIMARY KEY (id),
    CONSTRAINT owner_payout_items_payout_id_fkey
        FOREIGN KEY (payout_id) REFERENCES public.owner_payouts(id) ON DELETE CASCADE,
    CONSTRAINT owner_payout_items_settlement_id_fkey
        FOREIGN KEY (settlement_id) REFERENCES public.owner_settlements(id) ON DELETE CASCADE,
    CONSTRAINT owner_payout_items_unique UNIQUE (payout_id, settlement_id)
);

CREATE INDEX idx_owner_payout_items_payout_id
    ON public.owner_payout_items(payout_id);

CREATE INDEX idx_owner_payout_items_settlement_id
    ON public.owner_payout_items(settlement_id);

1.6 ตรวจสอบให้มี default account แค่ 1 บัญชีต่อ owner

อันนี้สำคัญมาก
CREATE UNIQUE INDEX uq_owner_bank_accounts_default_per_user
ON public.owner_bank_accounts(user_id)
WHERE is_default = true;

1.7 optional: trigger auto update updated_at

ถ้าคุณมี trigger กลางใช้อยู่แล้ว ข้ามได้
แต่ถ้ายังไม่มี ผมแนะนำทำ function กลาง

CREATE OR REPLACE FUNCTION public.set_updated_at()
RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_owner_bank_accounts_updated_at
BEFORE UPDATE ON public.owner_bank_accounts
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

CREATE TRIGGER trg_owner_settlements_updated_at
BEFORE UPDATE ON public.owner_settlements
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

CREATE TRIGGER trg_owner_payouts_updated_at
BEFORE UPDATE ON public.owner_payouts
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();