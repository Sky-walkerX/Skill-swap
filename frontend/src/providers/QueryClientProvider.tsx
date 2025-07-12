'use client';

import { QueryClient, QueryClientProvider as QueryProvider } from '@tanstack/react-query';
import { ReactNode, useState } from 'react';

export function QueryClientProvider({ children }: { children: ReactNode }) {
	const [client] = useState(() => new QueryClient());

	return <QueryProvider client={client}>{children}</QueryProvider>;
}
