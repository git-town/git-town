declare namespace delay {
	interface ClearablePromise<T> extends Promise<T> {
		/**
		Clears the delay and settles the promise.
		*/
		clear(): void;
	}

	/**
	Minimal subset of `AbortSignal` that delay will use if passed.
	This avoids a dependency on dom.d.ts.
	The dom.d.ts `AbortSignal` is compatible with this one.
	*/
	interface AbortSignal {
		readonly aborted: boolean;
		addEventListener(
			type: 'abort',
			listener: () => void,
			options?: {once?: boolean}
		): void;
		removeEventListener(type: 'abort', listener: () => void): void;
	}

	interface Options {
		/**
		An optional AbortSignal to abort the delay.
		If aborted, the Promise will be rejected with an AbortError.
		*/
		signal?: AbortSignal;
	}
}

type Delay = {
	/**
	Create a promise which resolves after the specified `milliseconds`.

	@param milliseconds - Milliseconds to delay the promise.
	@returns A promise which resolves after the specified `milliseconds`.
	*/
	(milliseconds: number, options?: delay.Options): delay.ClearablePromise<void>;

	/**
	Create a promise which resolves after the specified `milliseconds`.

	@param milliseconds - Milliseconds to delay the promise.
	@returns A promise which resolves after the specified `milliseconds`.
	*/
	<T>(
		milliseconds: number,
		options?: delay.Options & {
			/**
			Value to resolve in the returned promise.
			*/
			value: T;
		}
	): delay.ClearablePromise<T>;

	/**
	Create a promise which rejects after the specified `milliseconds`.

	@param milliseconds - Milliseconds to delay the promise.
	@returns A promise which rejects after the specified `milliseconds`.
	*/
	// TODO: Allow providing value type after https://github.com/Microsoft/TypeScript/issues/5413 will be resolved.
	reject(
		milliseconds: number,
		options?: delay.Options & {
			/**
			Value to reject in the returned promise.
			*/
			value?: unknown;
		}
	): delay.ClearablePromise<never>;
};

declare const delay: Delay & {
	createWithTimers(timers: {
		clearTimeout: typeof clearTimeout;
		setTimeout: typeof setTimeout;
	}): Delay;

	// TODO: Remove this for the next major release
	default: typeof delay;
};

export = delay;
