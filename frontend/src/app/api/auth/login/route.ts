import { NextResponse } from 'next/server'


export async function POST(req: Request) {
  const { email, password } = await req.json()


  return NextResponse.json({
    // id: '123', 
    // email,
    // name: 'abc'  dummy for now
  })
}
