import { BabylonOptions } from 'babylon';
import * as b from 'babel-types';
export { BabylonOptions };
export interface ExpressionToConstantOptions {
    constants?: any;
}
export interface Options extends ExpressionToConstantOptions {
    babylon?: BabylonOptions;
}
export declare function expressionToConstant(expression: b.Expression, options?: ExpressionToConstantOptions): {
    constant: true;
    result: any;
} | {
    constant: false;
    result?: void;
};
export declare function isConstant(src: string, constants?: any, options?: BabylonOptions): boolean;
export declare function toConstant(src: string, constants?: any, options?: BabylonOptions): any;
export default isConstant;
