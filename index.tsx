import React from 'react';
import {render} from 'ink';
import {App} from './src/App';
import chalk from 'chalk';

// Force level 3 (truecolor) with hex color support
chalk.level = 3;

void render(<App />);
