import React from 'react';
import { Box, Text } from 'ink';

interface PromptProps {
    value: string
}

const borderStyle = {
    top: '',
    bottom: '',
    left: '█',
    right: '',
    topLeft: '',
    topRight: '',
    bottomLeft: '',
    bottomRight: ''
}

export const Prompt: React.FC<PromptProps> = ({ value }: PromptProps) => {

    const lines = value.split('\n');

    return (
        <Box borderStyle={borderStyle} borderColor="#838ba7" padding={1} flexDirection="column">
            <Box flexDirection="column">
                {lines.length > 0 && (
                    lines.map((line: string, index: number) => (
                        <Text key={index}>
                            {line.trim()}
                        </Text>
                    ))
                )}
            </Box>
        </Box>
    );
};
