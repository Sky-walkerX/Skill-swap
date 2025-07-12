import { NextResponse } from 'next/server';
import { z } from 'zod';

const RegisterSchema = z.object({
	name: z.string().min(1),
	email: z.email(),
	password: z.string().min(6)
});

export async function POST(req: Request) {
	const body = await req.json();
	const parsed = RegisterSchema.safeParse(body);

	if (!parsed.success) {
		return NextResponse.json({ error: 'Invalid input' }, { status: 400 });
	}

	const { name, email, password } = parsed.data;

	return NextResponse.json({ success: true, user: { name, email } }, { status: 201 });
}
