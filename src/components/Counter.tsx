import React, {useState, useEffect} from 'react';
import {Text} from 'ink';

export const Counter: React.FC = () => {
	const [counter, setCounter] = useState(0);

	useEffect(() => {
		const timer = setInterval(() => {
			setCounter((previousCounter) => previousCounter + 1);
		}, 100);

		return () => {
			clearInterval(timer);
		};
	}, []);

	return <Text color="#266bbb">{counter} tests passed</Text>;
};
