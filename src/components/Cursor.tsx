import React, { useState, useEffect } from 'react';
import { Text } from 'ink';

interface CursorState {
	visible: boolean;
}

export const Cursor: React.FC = () => {
	const [cursorState, setCursorState] = useState<CursorState>({ visible: true });

	useEffect(() => {
		const interval = setInterval(() => {
			setCursorState((prev) => ({ visible: !prev.visible }));
		}, 500);

		return () => clearInterval(interval);
	}, []);

	return cursorState.visible ? (
		<Text backgroundColor="white" color="black">
			{' '}
		</Text>
	) : (
		<Text>{' '}</Text>
	);
};
