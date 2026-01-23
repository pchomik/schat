import React, { useState } from 'react';
import { Box, Text } from 'ink';
import { useInput } from 'ink';
import { Cursor } from './Cursor';


interface TextAreaProps {
    newPromptCallback: (newPrompt: string) => void
}


export const TextArea: React.FC<TextAreaProps> = React.memo(({newPromptCallback}: TextAreaProps) => {
	const [value, setValue] = useState('');

	useInput((input, key) => {
        if (key.ctrl && input === 's') {
            newPromptCallback(value);
            setValue('');
            return;
        } else if (key.return) {
            setValue((prev) => prev + '\n');
		} else if (key.backspace || key.delete) {
			setValue((prev) => prev.slice(0, -1));
		} else if (input && !key.ctrl && !key.meta) {
			setValue((prev) => prev + input);
		}
	});

	const lines = value.split('\n');

	return (
        <Box borderStyle="round" borderColor="#838ba7" padding={1} flexDirection="column">
            <Box flexDirection="column">
                {lines.length > 0 && (
                    lines.map((line, index) => (
                        <Text key={index}>
                            {line}
                            {index === lines.length - 1 && <Cursor />}
                        </Text>
                    ))
                )}
                {lines.length === 0 && <Cursor />}
            </Box>
        </Box>
	);
});
