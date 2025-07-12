'use client';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { SiGithub, SiGoogle } from 'react-icons/si';
import Image from 'next/image';
import { useState } from 'react';
import { signIn } from 'next-auth/react';
import { useRouter } from 'next/navigation';

interface UserDataType {
	email: string;
	password: string;
}

export function LoginForm({ className, ...props }: React.ComponentProps<'div'>) {
	const [userData, setUserData] = useState<UserDataType>({
		email: '',
		password: ''
	});
	const [error, setError] = useState('');
	const router = useRouter();
	const [loading, setLoading] = useState(false);

	const handleLogin = async (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		setError('');
		setLoading(true);

		const res = await signIn('credentials', {
			...userData,
			redirect: false
		});

		if (res?.error) {
			setError('Invalid credentials');
		} else {
			router.push('/'); // or wherever you want after login
		}

		setLoading(false);
	};

	return (
		<div className={cn('flex flex-col gap-6', className)} {...props}>
			<Card className="overflow-hidden p-0">
				<CardContent className="grid p-0 md:grid-cols-2">
					<form className="p-6 md:p-8" onSubmit={handleLogin}>
						<div className="flex flex-col gap-6">
							<div className="flex flex-col items-center text-center">
								<h1 className="text-2xl font-bold">Welcome back</h1>
								<p className="text-muted-foreground text-balance">
									Login to your SkillSwap account
								</p>
							</div>
							<div className="grid gap-3">
								<Label htmlFor="email">Email</Label>
								<Input
									id="email"
									type="email"
									placeholder="abc@example.com"
									required
									value={userData.email}
									onChange={(e) =>
										setUserData({ ...userData, email: e.target.value })
									}
								/>
							</div>
							<div className="grid gap-3">
								<div className="flex items-center">
									<Label htmlFor="password">Password</Label>
									<a
										href="#"
										className="ml-auto text-sm underline-offset-2 hover:underline"
									>
										Forgot your password?
									</a>
								</div>
								<Input
									id="password"
									type="password"
									required
									value={userData.password}
									onChange={(e) =>
										setUserData({ ...userData, password: e.target.value })
									}
									placeholder="Enter your password"
								/>
							</div>
							{error && <p className="text-red-500">{error}</p>}
							<Button
								type="submit"
								className="w-full cursor-pointer"
								disabled={loading}
							>
								Login
							</Button>
							<div className="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t">
								<span className="bg-card text-muted-foreground relative z-10 px-2">
									Or continue with
								</span>
							</div>
							<div className="grid grid-cols-2 gap-4">
								<Button
									variant="outline"
									type="button"
									className="w-full"
									disabled={loading}
									onClick={() => signIn('github')}
								>
									<span className="flex gap-2 justify-center items-center">
										<SiGithub className="h-4 w-4" />
										Github
									</span>
									<span className="sr-only">Login with Github</span>
								</Button>
								<Button
									variant="outline"
									type="button"
									className="w-full"
									disabled={loading}
									onClick={() => signIn('github')}
								>
									<span className="flex gap-2 justify-center items-center">
										<SiGoogle className="h-4 w-4" />
										Google
									</span>
									<span className="sr-only">Login with Google</span>
								</Button>
							</div>
							<div className="text-center text-sm">
								Don&apos;t have an account?{' '}
								<a href="#" className="underline underline-offset-4">
									Sign up
								</a>
							</div>
						</div>
					</form>
					<div className="bg-muted relative hidden md:block">
						<Image
							src="https://media.gettyimages.com/id/1125868664/video/making-smart-moves-across-the-digital-landscape.jpg?s=640x640&k=20&c=94zKlF9shOT1fiQXShCgJHB2X2_AkYAjTJ4tsj-3uTs="
							alt="Image"
							width={500}
							height={500}
							className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
						/>
					</div>
				</CardContent>
			</Card>
			<div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
				By clicking continue, you agree to our <a href="#">Terms of Service</a> and{' '}
				<a href="#">Privacy Policy</a>.
			</div>
		</div>
	);
}
