import GameplayKit
import SpriteKit

class GameScene: SKScene {
    var player = SKSpriteNode()
    var opponent = SKSpriteNode()
    var ball = SKSpriteNode()

    override func didMove(to view: SKView) {
        player = self.childNode(withName: "player") as! SKSpriteNode
        opponent = self.childNode(withName: "opponent") as! SKSpriteNode
        ball = self.childNode(withName: "ball") as! SKSpriteNode

        ball.physicsBody?.applyImpulse(CGVector(dx: 40, dy: 40))

        let border = SKPhysicsBody(edgeLoopFrom: self.frame)
        border.friction = 0
        border.restitution = 1

        self.physicsBody = border
    }

    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        for touch in touches {
            let location = touch.location(in: self)
            player.run(SKAction.moveTo(x: location.x, duration: 0.0))
        }
    }

    override func touchesMoved(_ touches: Set<UITouch>, with event: UIEvent?) {
        for touch in touches {
            let location = touch.location(in: self)
            player.run(SKAction.moveTo(x: location.x, duration: 0.0))
        }
    }

    override func update(_ currentTime: TimeInterval) {
        opponent.run(SKAction.moveTo(x: ball.position.x, duration: 0.0))
    }
}
